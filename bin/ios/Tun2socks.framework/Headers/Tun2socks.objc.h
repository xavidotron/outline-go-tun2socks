// Objective-C API for talking to github.com/Jigsaw-Code/outline-go-tun2socks/ios Go package.
//   gobind -lang=objc github.com/Jigsaw-Code/outline-go-tun2socks/ios
//
// File is generated by gobind. Do not edit.

#ifndef __Tun2socks_H__
#define __Tun2socks_H__

@import Foundation;
#include "Universe.objc.h"


@protocol Tun2socksPacketFlow;
@class Tun2socksPacketFlow;

@protocol Tun2socksPacketFlow <NSObject>
- (void)writePacket:(NSData*)packet;
@end

FOUNDATION_EXPORT void Tun2socksInputPacket(NSData* data);

FOUNDATION_EXPORT void Tun2socksStartSocks(id<Tun2socksPacketFlow> packetFlow, NSString* proxyHost, long proxyPort);

FOUNDATION_EXPORT void Tun2socksStopSocks(void);

@class Tun2socksPacketFlow;

@interface Tun2socksPacketFlow : NSObject <goSeqRefInterface, Tun2socksPacketFlow> {
}
@property(strong, readonly) id _ref;

- (instancetype)initWithRef:(id)ref;
- (void)writePacket:(NSData*)packet;
@end

#endif
